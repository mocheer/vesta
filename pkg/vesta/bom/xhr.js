// 保存原始的 XMLHttpRequest 构造函数
const OriginalXMLHttpRequest = window.XMLHttpRequest;

// 创建一个新的 XMLHttpRequest 构造函数
function CustomXMLHttpRequest() {
  const xhr = new OriginalXMLHttpRequest();

  // 拦截发送请求
  const originalSend = xhr.send;
  xhr.send = function() {
    // 可以在这里添加自定义逻辑，比如修改请求参数、添加额外的请求头等
    console.log('Intercepted request:', this);

    // 调用原始的 send 方法
    originalSend.apply(xhr, arguments);
  };

  // 拦截 open 请求
  const originalOpen = xhr.open;
  xhr.open = function(method, url, async, user, password) {
    // 可以在这里添加自定义逻辑，比如修改请求 URL 等
    console.log('Intercepted open:', method, url);

    // 调用原始的 open 方法
    originalOpen.apply(xhr, arguments);
  };

  return xhr;
}

// 重写全局的 XMLHttpRequest 构造函数
window.XMLHttpRequest = CustomXMLHttpRequest;