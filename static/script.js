document.getElementById("comprimirBtn").onclick = () => {
    const file = document.getElementById("archivo").files[0];
    if (!file) return alert("Selecciona un archivo.");

    // Simula compresión
    const orig = file.size;
    const comp = Math.floor(orig * 0.5);
    const ratio = (comp / orig).toFixed(2);

    document.getElementById("interfaz1").classList.add("hidden");
    document.getElementById("interfaz2").classList.remove("hidden");

    document.getElementById("original").innerText = orig;
    document.getElementById("comprimido").innerText = comp;
    document.getElementById("ratio").innerText = ratio;
};

document.getElementById("descargarBtn").onclick = () => {
    alert("Aquí descarga real o redirección a /descargar");
};
