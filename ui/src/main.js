import App from './App.svelte';

let div = document.createElement('div');
let script = document.currentScript;
script.parentNode.insertBefore(div, script);

const app = new App({
	target: div,
});

export default app;