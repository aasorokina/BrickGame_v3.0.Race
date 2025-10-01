import { GAME_FIELD_WIDTH, GAME_FIELD_HEIGHT, tileColors } from './config.js';

export class GameField {
    constructor($gameField) {
        this.element = $gameField;
        this.tiles = [];
        for (let y = 0; y < GAME_FIELD_HEIGHT; ++y) {
            for (let x = 0; x < GAME_FIELD_WIDTH; ++x) {
                const $tile = document.createElement('div');
                $tile.classList.add('tile');
                $tile.id = `position-${x}-${y}`;
                this.tiles.push($tile);
                this.element.append($tile);
            }
        }
    }


    getTile(x, y) {
        return this.tiles[y * GAME_FIELD_WIDTH + x];
    }

    enableTile(x, y, color) {
        this.getTile(x, y).style.background = color;
    }

    disableTile(x, y) {
        this.getTile(x, y).style.background = tileColors[0];
    }

    disableAllTiles() {
        this.tiles.forEach((tile) => {
            tile.style.background = tileColors[0];
        })
    }
}
