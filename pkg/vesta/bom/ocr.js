// import { createWorker } from 'tesseract.js';
// https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.esm.min.js
// https://github.com/naptha/tesseract.js
import { createWorker } from 'https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.esm.min.js';

(async () => {
  const worker = await createWorker('eng');
  // const ret = await worker.recognize('https://tesseract.projectnaptha.com/img/eng_bw.png');
  const t = new Date().toISOString().slice(0,-1)
  const ret = await worker.recognize("http://10.135.6.98/Account/GetValidateCode?time="+t+"&browserName=Safari")
  console.log(ret.data.text);
  await worker.terminate();
})();


function loadAndRunTesseract() {
    // 1. 创建type="module"的script标签
    const script = document.createElement('script');
    script.type = 'module';
    // 2. 设置script内容（内联模块代码）
    script.innerHTML = `
import { createWorker } from 'https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.esm.min.js';

(async () => {
  const worker = await createWorker('eng');
  // const ret = await worker.recognize('https://tesseract.projectnaptha.com/img/eng_bw.png');
  const t = new Date().toISOString().slice(0,-1)
  const ret = await worker.recognize("http://10.135.6.98/Account/GetValidateCode?time="+t+"&browserName=Safari")
  console.log(ret.data.text);
  await worker.terminate();
})();
    `;
    
    // 3. 添加到文档中执行
    document.head.appendChild(script);
}