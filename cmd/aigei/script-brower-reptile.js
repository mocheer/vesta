async function reptile() {
  let lis = Array.from(document.querySelector("#explorerTree_74294").children);
  let extendSpans = Array.from(lis).map((e) => e.children[0]);
  for (let i = 0; i < lis.length; i++) {
    let span = extendSpans[i];
    let ulChildren = lis[i].querySelector("ul")?.children;
    if (!ulChildren) {
      span.click(); // 点击展开
      await new Promise((resolve) => {
        setTimeout(() => {
          ulChildren = lis[i].querySelector("ul")?.children;
          resolve();
        }, 3000);
      });
    }
    if (ulChildren.length > 0) {
      for (let j = 0; j < ulChildren.length; j++) {
        ulChildren[j].children[1].click();
        await new Promise((resolve) => {
          setTimeout(() => {
            resolve();
          }, 2000);
        });
      }
    }
  }
  return "success";
}
await reptile();
