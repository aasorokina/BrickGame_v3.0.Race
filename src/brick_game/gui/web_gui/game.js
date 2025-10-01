import { applyRootStyles } from './src/utils.js';
import { GameField } from './src/game-field.js';
import { rootStyles, keyCodes } from './src/config.js';
import { getGamesList } from './src/games-list.js';
import { selectGame } from './src/games-list.js';
import { NextField } from './src/next-field.js';
import { setAction } from './src/user-input.js';
import { getCurrentState } from './src/update-current-state.js';

applyRootStyles(rootStyles);
const gameField = new GameField(document.querySelector('#game-field'));
const nextField = new NextField(document.querySelector('#next-field'));
let pressedKey = -1;

const $rightPanel = document.querySelector('#right-panel');
const leftPanel = document.querySelector('#left-panel');

const notification = document.createElement('div');
notification.id = 'notification';
notification.textContent = '';
leftPanel.appendChild(notification);

document.addEventListener('DOMContentLoaded', getGamesList);
document.addEventListener('keydown', async function (event) {
    if (keyCodes.game_1.includes(event.code)) {
        await selectGame(1);
    } else if (keyCodes.game_2.includes(event.code)) {
        await selectGame(2);
    } else if (keyCodes.game_3.includes(event.code)) {
        await selectGame(3);
    } else if (keyCodes.up.includes(event.code)) {
        pressedKey = 5;
    } else if (keyCodes.right.includes(event.code)) {
        pressedKey = 4;
    } else if (keyCodes.down.includes(event.code)) {
        pressedKey = 6;
    } else if (keyCodes.left.includes(event.code)) {
        pressedKey = 3;
    } else if (keyCodes.start.includes(event.code)) {
        pressedKey = 0;
    } else if (keyCodes.pause.includes(event.code)) {
        pressedKey = 1;
    } else if (keyCodes.action.includes(event.code)) {
        pressedKey = 7;
    } else if (keyCodes.exit.includes(event.code)) {
        gameField.disableAllTiles();
        nextField.disableAllTiles();
        document.querySelectorAll('#games-list button')
            .forEach(btn => btn.classList.remove('active'));
        document.querySelector('#score-value').textContent = '0';
        document.querySelector('#high-score-value').textContent = '0';
        document.querySelector('#level-value').textContent = '0';
        document.querySelector('#speed-value').textContent = '0';
        pressedKey = 2;
    }
});

let isGameRunning = true;

function gameLoop() {
    if (!isGameRunning) return;

    if (pressedKey >= 0) {
        setAction(pressedKey, false);
        pressedKey = -1;
    }
    getCurrentState(gameField, nextField, $rightPanel, notification);
    requestAnimationFrame(gameLoop);
}

requestAnimationFrame(gameLoop);