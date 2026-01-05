let currentOrderId = null;

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('createOrderBtn').addEventListener('click', createOrder);
    document.getElementById('refreshBtn').addEventListener('click', refreshOrderStatus);
    document.getElementById('syncBtn').addEventListener('click', syncOrderStatus);
    document.getElementById('refreshOrdersBtn').addEventListener('click', loadOrders);
    document.getElementById('ordersTable').addEventListener('click', onOrdersTableClick);
    
    loadOrders();
});

async function createOrder() {
    const amount = document.getElementById('amount').value;
    const subject = document.getElementById('subject').value;
    const body = document.getElementById('body').value;
    const payType = document.getElementById('payType').value;

    if (!amount || !subject) {
        alert('请填写支付金额和商品名称');
        return;
    }

    try {
        let endpoint, requestData;
        
        if (payType === 'app') {
            endpoint = '/api/app-orders';
            requestData = {
                total_amount: amount,
                subject: subject,
                body: body
            };
        } else {
            endpoint = '/api/orders';
            requestData = {
                total_amount: amount,
                subject: subject,
                body: body
            };
        }

        const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        });

        if (!response.ok) {
            throw new Error('创建订单失败');
        }

        const data = await response.json();
        currentOrderId = data.order_id;

        if (payType === 'app') {
            displayAppPayment(data.pay_url, data.order_id, amount);
        } else {
            displayQRCode(data.qr_code, data.order_id, amount);
        }
        
        startPolling();

    } catch (error) {
        console.error('Error:', error);
        alert('创建订单失败: ' + error.message);
    }
}

function displayQRCode(qrCode, orderId, amount) {
    const container = document.getElementById('qrCodeContainer');
    const qrCodeDiv = document.getElementById('qrCode');
    const orderIdSpan = document.getElementById('orderId');
    const orderAmountSpan = document.getElementById('orderAmount');

    qrCodeDiv.innerHTML = `<img src="${qrCode}" alt="支付二维码">`;
    orderIdSpan.textContent = orderId;
    orderAmountSpan.textContent = amount;
    container.style.display = 'block';

    document.getElementById('qrCodeContainer').scrollIntoView({ behavior: 'smooth' });
}

function displayAppPayment(payUrl, orderId, amount) {
    const container = document.getElementById('qrCodeContainer');
    const qrCodeDiv = document.getElementById('qrCode');
    const orderIdSpan = document.getElementById('orderId');
    const orderAmountSpan = document.getElementById('orderAmount');

    const safeUrl = payUrl || '';
    qrCodeDiv.innerHTML = `
        <div style="text-align: center;">
            <p>正在跳转到支付宝...</p>
            <p>如果页面没有自动跳转，请点击下方按钮</p>
            <a href="${safeUrl}" class="btn btn-primary" target="_blank" rel="noreferrer">立即支付</a>
        </div>
    `;
    orderIdSpan.textContent = orderId;
    orderAmountSpan.textContent = amount;
    container.style.display = 'block';

    setTimeout(() => {
        if (safeUrl) {
            window.location.href = safeUrl;
        }
    }, 1000);
}

async function refreshOrderStatus() {
    if (!currentOrderId) {
        alert('请先创建订单');
        return;
    }

    try {
        const response = await fetch(`/api/orders/${currentOrderId}`);
        if (!response.ok) {
            throw new Error('获取订单状态失败');
        }

        const order = await response.json();
        displayOrderStatus(order);

    } catch (error) {
        console.error('Error:', error);
        alert('获取订单状态失败: ' + error.message);
    }
}

async function syncOrderStatus() {
    if (!currentOrderId) {
        alert('请先创建订单');
        return;
    }

    const order = await syncOrderStatusById(currentOrderId);
    if (order) {
        displayOrderStatus(order);
    }
}

async function syncOrderStatusById(orderId) {
    try {
        const response = await fetch(`/api/orders/${orderId}/sync`, { method: 'POST' });
        if (!response.ok) {
            const data = await response.json().catch(() => null);
            const message = data && (data.detail || data.error) ? (data.detail || data.error) : '更新支付状态失败';
            throw new Error(message);
        }

        const data = await response.json();
        const order = data.order;
        order.alipay_trade_status = data.alipay_trade_status;
        loadOrders();
        return order;

    } catch (error) {
        console.error('Error:', error);
        alert('更新支付状态失败: ' + error.message);
        return null;
    }
}

async function onOrdersTableClick(e) {
    const target = e.target;
    if (!(target instanceof HTMLElement)) {
        return;
    }
    if (!target.classList.contains('sync-order-btn')) {
        return;
    }

    const orderId = target.dataset.orderId;
    if (!orderId) {
        return;
    }

    const originalText = target.textContent;
    target.textContent = '更新中...';
    target.setAttribute('disabled', 'disabled');

    try {
        const order = await syncOrderStatusById(orderId);
        if (order) {
            currentOrderId = orderId;
            displayOrderStatus(order);
        }
    } finally {
        target.textContent = originalText;
        target.removeAttribute('disabled');
    }
}

function displayOrderStatus(order) {
    const statusDiv = document.getElementById('orderStatus');
    
    let statusText = '';
    let statusClass = '';

    switch(order.status) {
        case 'pending':
            statusText = '等待支付';
            statusClass = 'pending';
            break;
        case 'paid':
            statusText = '支付成功';
            statusClass = 'success';
            break;
        case 'failed':
            statusText = '支付失败';
            statusClass = 'failed';
            break;
        case 'closed':
            statusText = '订单已关闭';
            statusClass = 'closed';
            break;
    }

    statusDiv.innerHTML = `
        <h3>订单状态</h3>
        <p class="${statusClass}">${statusText}</p>
        ${order.alipay_trade_status ? `<p>支付宝状态: ${order.alipay_trade_status}</p>` : ''}
        ${order.trade_no ? `<p>支付宝交易号: ${order.trade_no}</p>` : ''}
    `;
    statusDiv.style.display = 'block';

    if (order.status === 'paid') {
        stopPolling();
        loadOrders();
    }
}

let pollingInterval = null;

function startPolling() {
    stopPolling();
    pollingInterval = setInterval(refreshOrderStatus, 3000);
}

function stopPolling() {
    if (pollingInterval) {
        clearInterval(pollingInterval);
        pollingInterval = null;
    }
}

async function loadOrders() {
    try {
        const response = await fetch('/api/orders');
        if (!response.ok) {
            throw new Error('获取订单列表失败');
        }

        const orders = await response.json();
        displayOrders(orders);

    } catch (error) {
        console.error('Error:', error);
        alert('获取订单列表失败: ' + error.message);
    }
}

function displayOrders(orders) {
    const tableDiv = document.getElementById('ordersTable');

    if (orders.length === 0) {
        tableDiv.innerHTML = '<p style="text-align: center; color: #999;">暂无订单</p>';
        return;
    }

    const sortedOrders = orders.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

    tableDiv.innerHTML = `
        <table>
            <thead>
                <tr>
                    <th>订单号</th>
                    <th>商品名称</th>
                    <th>金额</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                ${sortedOrders.map(order => `
                    <tr>
                        <td>${order.id}</td>
                        <td>${order.subject}</td>
                        <td>¥${order.total_amount}</td>
                        <td><span class="status-badge ${order.status}">${getStatusText(order.status)}</span></td>
                        <td>${formatDate(order.created_at)}</td>
                        <td><button class="btn btn-secondary sync-order-btn" data-order-id="${order.id}">更新状态</button></td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

function getStatusText(status) {
    switch(status) {
        case 'pending':
            return '等待支付';
        case 'paid':
            return '支付成功';
        case 'failed':
            return '支付失败';
        case 'closed':
            return '订单已关闭';
        default:
            return status;
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
}
