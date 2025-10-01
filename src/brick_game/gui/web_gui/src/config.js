export const GAME_FIELD_WIDTH = 10;
export const GAME_FIELD_HEIGHT = 20;

export const NEXT_FIELD_WIDTH = 4;
export const NEXT_FIELD_HEIGHT = 4;

export const rootStyles = {
    '--tile-size': '20px',
    '--tile-color': '#242727ff',
    '--tile-active-color': '#222',
    '--game-field-width': GAME_FIELD_WIDTH,
    '--game-field-height': GAME_FIELD_HEIGHT,
    '--game-field-gap': '2px',
    '--game-field-background': '#333',
    '--next-field-width': NEXT_FIELD_WIDTH,
    '--next-field-height': NEXT_FIELD_HEIGHT,
    '--next-field-gap': '2px',
    '--next-field-background': '#333',
};

export const keyCodes = {
    up: ['ArrowUp'],
    right: ['ArrowRight'],
    down: ['ArrowDown'],
    left: ['ArrowLeft'],
    pause: ['KeyP', 'Keyp'],
    exit: ['KeyQ', 'Keyq'],
    action: ['Space'],
    start: ['Enter'],
    game_1: ['Digit1', 'Numpad1'],
    game_2: ['Digit2', 'Numpad2'],
    game_3: ['Digit3', 'Numpad3'],
};

export const tileColors = {
    0: '#242727ff',
    1: '#d6204e',
    2: '#0b903f',
    4: '#0066cc',
    6: '#30a7cbff',
    8: '#ffa505',
    9: '#fffb05',
    10: '#6900a3',
    11: '#459f69',
};
