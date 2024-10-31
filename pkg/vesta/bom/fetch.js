// 保存原始的 fetch 函数
const originalFetch = window.fetch;

// 自定义的 fetch 函数
function customFetch(url, options) {
  return new Promise((resolve, reject) => {
    // 拦截请求
    console.log('Intercepted URL:', url);

    // 可以在这里添加自定义逻辑，比如修改请求参数、添加额外的请求头等

    // 调用原始的 fetch 函数
    originalFetch(url, options)
      .then(response => {
        // 可以在这里添加对响应的拦截逻辑
        console.log('Response received:', response);
        resolve(response);
      })
      .catch(error => {
        reject(error);
      });
  });
}

// 重写全局的 fetch 函数
window.fetch = customFetch;