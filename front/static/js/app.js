let currentOrderId = null;

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('createOrderBtn').addEventListener('click', createOrder);
    document.getElementById('refreshBtn').addEventListener('click', refreshOrderStatus);
    document.getElementById('refreshOrdersBtn').addEventListener('click', loadOrders);
    
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
            displayAppPayment(data.order_str, data.order_id, amount);
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

function displayAppPayment(orderStr, orderId, amount) {
    const container = document.getElementById('qrCodeContainer');
    const qrCodeDiv = document.getElementById('qrCode');
    const orderIdSpan = document.getElementById('orderId');
    const orderAmountSpan = document.getElementById('orderAmount');

    // 创建一个隐藏的表单来提交支付请求
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = 'https://openapi-sandbox.dl.alipaydev.com/gateway.do'; // 生产环境
    // form.action = 'https://openapi.alipaydev.com/gateway.do'; // 沙箱环境
    form.style.display = 'none';

    // 解析 orderStr 为参数对
    const params = orderStr.split('&');
    params.forEach(param => {
        const [key, value] = param.split('=');
        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = decodeURIComponent(key);
        input.value = decodeURIComponent(value);
        form.appendChild(input);
    });

    document.body.appendChild(form);

    // 添加提示信息
    qrCodeDiv.innerHTML = `
        <div style="text-align: center;">
            <p>正在跳转到支付宝...</p>
            <p>如果页面没有自动跳转，请点击下方按钮</p>
            <button onclick="document.querySelector('form').submit();" class="btn btn-primary">立即支付</button>
        </div>
    `;
    orderIdSpan.textContent = orderId;
    orderAmountSpan.textContent = amount;
    container.style.display = 'block';

    // 自动提交表单
    setTimeout(() => {
        form.submit();
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
