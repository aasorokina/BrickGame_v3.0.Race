import { tileColors } from './config.js';

export async function getCurrentState(gameField, nextField, $rightPanel, notification) {
    try {
        const response = await fetch('/api/state');
        if (response.ok) {
            const state = await response.json();
            console.log(state);
            render(gameField, nextField, $rightPanel, state);
        } else {
            notification.textContent = 'Choose the game';
            console.error(response.statusText);
        }
    } catch (error) {
        console.error(error);
    }
}

function render(gameField, nextField, $rightPanel, state) {
    $rightPanel.querySelector('#score-value').textContent = state.score;
    $rightPanel.querySelector('#high-score-value').textContent = state.high_score;
    $rightPanel.querySelector('#level-value').textContent = state.level;
    $rightPanel.querySelector('#speed-value').textContent = state.speed;
    gameField.disableAllTiles()
    state.field.forEach((row, i) => {
        row.forEach((elem, j) => {
            gameField.enableTile(j, i, tileColors[elem]);
        });
    });
    nextField.disableAllTiles();
    state.next.forEach((row, i) => {
        row.forEach((elem, j) => {
            nextField.enableTile(j, i, tileColors[elem]);
        });
    });

    if (state.pause == -1) {
        notification.textContent = 'Press Enter to start';
    } else if (state.pause == 1) {
        notification.textContent = '[Paused]';
    } else if (state.pause >= 2) {
        notification.textContent = 'GameOver! \n Press Enter to restart \n or Q to exit';
    } else {
        notification.textContent = '';
    }
}