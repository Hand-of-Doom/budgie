import {render, html} from 'lit-html';

function home() {
    return html`<h1>Hello world!</h1>`;
}

render(home(), document.body);
