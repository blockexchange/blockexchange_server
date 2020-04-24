import state from './state.js';

export default {
  view: function(){
    return m("div", { class: "progress"}, [
      m("div", { class: "progress-bar", style: `width: ${state.progress}%` }, [
        (Math.floor(state.progress * 10) / 10) + "%"
      ])
    ]);
  }
};
