let processedFileName = '';

async function handleCompression(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const fileInput = event.target.querySelector('input[type="file"]');
    if (!fileInput.files[0]) {
        alert('Por favor seleccione un archivo');
        return;
    }
    formData.append('fileName', fileInput.files[0].name);
    
    try {
        const response = await fetch('/compress', {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) throw new Error('Compression failed');
        processedFileName = await response.text();
        document.getElementById('downloadBtn').style.display = 'block';
    } catch (error) {
        alert('Error: ' + error.message);
    }
}

async function handleDecompression(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const fileInput = event.target.querySelector('input[type="file"]');
    if (!fileInput.files[0]) {
        alert('Por favor seleccione un archivo');
        return;
    }
    formData.append('fileName', fileInput.files[0].name);
    
    try {
        const response = await fetch('/decompress', {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) throw new Error('Decompression failed');
        processedFileName = await response.text();
        document.getElementById('downloadBtn').style.display = 'block';
    } catch (error) {
        alert('Error: ' + error.message);
    }
}

function downloadFile() {
    if (processedFileName) {
        window.location.href = `/download?file=${encodeURIComponent(processedFileName)}`;
    }
}

// Asignar los event listeners
document.getElementById('compressForm')?.addEventListener('submit', handleCompression);
document.getElementById('decompressForm')?.addEventListener('submit', handleDecompression);
document.getElementById('downloadBtn')?.addEventListener('click', downloadFile);
