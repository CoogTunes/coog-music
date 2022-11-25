function tabControl(){
  document.addEventListener('click', function(evt) {
    let target = evt.target;
    if(target.matches('.tab-item')){
      let targetTab = target.getAttribute('data-container-target');
      let parentElement = target.parentElement;
      toggleTab(target, targetTab, parentElement);
    }
  });

  function toggleTab(target, targetTab, parent){
    let parentContainer = parent.parentElement;
    currentlySelected(parentContainer);
    let selectedTab = parentContainer.querySelector(targetTab);
    selectedTab.classList.toggle('selected');
    target.classList.toggle('selected');
  }

  function currentlySelected(mainParent){
    let currentlySelected = mainParent.querySelector('.tab-item.selected');
    let currentTabSelected = mainParent.querySelector(currentlySelected.getAttribute('data-container-target'));
    currentlySelected.classList.toggle('selected');
    currentTabSelected.classList.toggle('selected');
  }

}

window.addEventListener("DOMContentLoaded", function (evt) {
  tabControl();
});