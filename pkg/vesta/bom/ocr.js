import { createWorker } from 'tesseract.js';
// https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.min.js

(async () => {
  const worker = await createWorker('eng');
  const ret = await worker.recognize('https://tesseract.projectnaptha.com/img/eng_bw.png');
  console.log(ret.data.text);
  await worker.terminate();
})();