export async function getGamesList() {
    try {
        const response = await fetch('/api/games');
        if (response.ok) {
            const data = await response.json();
            setGameButtons(data.games);
        } else {
            console.error(response.statusText);
        }
    } catch (error) {
        console.error(error);
    }
}

function setGameButtons(games) {
    const gamesListContainer = document.getElementById('games-list');
    gamesListContainer.innerHTML = '';

    games.forEach(game => {
        const button = document.createElement('button');
        button.textContent = `${game.id}. ${game.name}`;
        button.onclick = () => selectGame(game.id);
        gamesListContainer.appendChild(button);
    });
}

export async function selectGame(gameId) {
    const buttons = document.querySelectorAll('#games-list button');
    try {
        const response = await fetch(`/api/games/${gameId}`, {
            method: 'POST'
        });

        if (response.ok) {
            buttons.forEach(btn => btn.classList.remove('active'));
            buttons[gameId - 1].classList.add('active');
            return true;
        } else {
            const errorData = await response.json().catch(() => ({}));
            console.error(errorData.message || 'Ошибка при выборе игры');
            return false;
        }
    } catch (err) {
        console.error(err);
        return false;
    }
}