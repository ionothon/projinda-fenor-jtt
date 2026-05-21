async function fetchImages(filterParam) {
    
    const url = `/api/images?${filterParam}`;

    try {

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

const dropZone = document.getElementById('drop-zone');
const fileInput = document.getElementById('file-input');

dropZone.addEventListener('click', () => fileInput.click());
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('hover');
});

dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('hover');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('hover');

    const files = e.dataTransfer.files;
    if (files.length > 0) {
        uploadFiles(files);
    }
});

async function uploadFiles(files) {
    const formData = new FormData();

    for (let i = 0; i < files.length; i++) {
        formData.append('uploadedImages', files[i]);
    }
    try {
        const response = await fetch('/api/upload', {
            method: 'POST',
            body: formData
        });
        if (response.ok) {
            fetchImages('');
        } else {
            alert('Upload failed');
        }
    } catch (error) {
        console.error('Could not load', error);
        alert('Could not reach server');
    }
}

function filterChanged() {
   const checkedColors = Array.from(document.querySelectorAll('.color-filter:checked'))
   .map(cb => cb.value);

   const colorParam = checkedColors.join(',');
   fetchImages(`color=${colorParam}`);
}

function clearColorFilters() {
    const colorCheckboxes = document.querySelectorAll('.color-filter');
    colorCheckboxes.forEach(chckbox => {
        chckbox.checked = false;
    });
    filterChanged();
}

const firstButton = document.querySelector(".filter-selection button");
if (firstButton) {
    firstButton.classList.add("active");
}
document.querySelectorAll("filter-selection button").forEach(button => {
    button.addEventListener("click", function() {
        document.querySelectorAll("filter-selection button").forEach(btn => {
            btn.classList.remove("active");
        });
        this.classList.add("active");
    });
});

const script = document.createElement('script');
script.src = 'https://cdnjs.cloudflare.com/ajax/libs/jszip/3.10.1/jszip.min.js';
document.head.appendChild(script);

script.onload = () => {
    const downloadBtn = document.getElementById('download-zip-btn');
    if (!downloadBtn) return;
    downloadBtn.addEventListener('click', async () => {
        const zip = new JSZip();
        const images = document.querySelectorAll('.card img');

        if (images.length == 0) {
            alert('No pictures to download');
            return;
        }
        downloadBtn.innerText = 'Packing images...';

        for (let i = 0; i < images.length; i++) {
            try {
                const imgUrl = images[i].src;
                const response = await fetch(imgUrl);
                const blob = await response.blob();
                zip.file(`pic-${i + 1}.jpg`, blob);
            } catch (err) {
                console.error('Could not pack image:', err);
            }
        }
        zip.generateAsync({ type: 'blob' }).then((content) => {
            const element = document.createElement('a');
            element.href = URL.createObjectURL(content);
            element.download = 'my-color-analyzis.zip';
            document.body.appendChild(element);
            element.click();
            document.body.removeChild(element);
            downloadBtn.innerText = 'Download filtered images (.zip)';
        });
    });
};

fetchImages('');