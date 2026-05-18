async function fetchImages(filterParam) {
    
    try {
        const url = `/api/images?${filterParam}`;

        const response = await fetch(url);
        const images = await response.json();
        const container = document.getElementById('image-container');

        container.innerHTML = '';
        if (images.length == 0) {
            container.innerHTML = '<p>No pictures found.</p>';
            return;
        }
        images.forEach(image => {
            const cardHTML = `
            <div class="card">
                <img src="img/${image.name}" alt="${image.name}">
                <h4>${image.name}</h4>
                <p><strong>Color:</strong> ${image.category}</p>
                <small>
                    H: ${Math.round(image.hue)}° |
                    S: ${Math.round(image.sat * 100)}% |
                    V: ${Math.round(image.val * 100)}%
                </small>
            </div>
            `;
            container.innerHTML += cardHTML;
        });

    } catch (error) {
        console.error("Can not get images:", error);
        document.getElementById('image-container').innerHTML = `
            <p style="color:red; font-weight:bold;">
                Can't reach Server.
            </p>
            `;
        }
}

fetchImages('');