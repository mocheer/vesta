async function thumbnailReptile() {
  const el = document.querySelector(
    "#tab-mount-content > div.tab-mount-content-body-container.tab-pane.active > div.rescListPageNav.pager-big > div.blueFoot.pageNav.js-page-nav span.js-page-btn-right",
  );
  el.click();
  await new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, 9000);
  });
  thumbnailReptile();
}

thumbnailReptile();
