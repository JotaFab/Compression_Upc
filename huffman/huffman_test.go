package huffman

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHuffmanEncodeDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen bool // si queremos verificar que la longitud comprimida sea menor
	}{
		{
			name:    "texto simple",
			input:   "hello world",
			wantLen: true,
		},
		{
			name:    "texto repetitivo",
			input:   "aaaaabbbcc",
			wantLen: true,
		},
		{
			name:    "caracteres especiales",
			input:   "¡Hola, mundo! 123 🌎",
			wantLen: false, // con caracteres especiales no garantizamos compresión
		},
		{
			name:    "texto largo",
			input:   strings.Repeat("Lorem ipsum dolor sit amet ", 100),
			wantLen: true,
		},
		{
			name:    "un solo caracter",
			input:   "aaaaaaaaaa",
			wantLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Codificar
			encoded, codeMap := Encode(tt.input)

			// Verificar que el mapa de códigos no esté vacío
			if len(codeMap) == 0 {
				t.Error("codeMap está vacío")
			}

			// Verificar que no hay códigos duplicados
			usedCodes := make(map[string]bool)
			for _, code := range codeMap {
				if usedCodes[code] {
					t.Errorf("código duplicado encontrado: %s", code)
				}
				usedCodes[code] = true
			}

			// Calcular y mostrar estadísticas de compresión
			originalBits := len(tt.input) * 8
			compressedBits := len(encoded)
			ratio := float64(compressedBits) / float64(originalBits) * 100
			t.Logf("Estadísticas para %s:", tt.name)
			t.Logf("  - Texto original: %q", tt.input)
			t.Logf("  - Texto comprimido: %q", encoded)
			// Mostrar el mapa de códigos
			t.Logf("  - Mapa de códigos: %v", codeMap)
			t.Logf("  - Tamaño original: %d bits", originalBits)
			t.Logf("  - Tamaño comprimido: %d bits", compressedBits)
			t.Logf("  - Ratio de compresión: %.2f%%", ratio)

			// Verificar compresión solo si se espera
			if tt.wantLen && compressedBits >= originalBits {
				t.Errorf("no se logró la compresión esperada: original %d bits, comprimido %d bits",
					originalBits, compressedBits)
			}

			// Decodificar y verificar
			decoded := Decode(encoded, codeMap)
			if decoded != tt.input {
				t.Errorf("la decodificación no coincide con el original:\nquiero: %s\nobtenido: %s",
					tt.input, decoded)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "texto vacío",
			input: "",
		},
		{
			name:  "un solo caracter",
			input: "a",
		},
		{
			name:  "dos caracteres iguales",
			input: "aa",
		},
		{
			name:  "espacios en blanco",
			input: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, codeMap := Encode(tt.input)
			decoded := Decode(encoded, codeMap)
			if decoded != tt.input {
				t.Errorf("fallo en caso borde: %s\nquiero: %q\nobtenido: %q",
					tt.name, tt.input, decoded)
			}
		})
	}
}

func TestFileCompression(t *testing.T) {
	// Definir archivos de prueba
	testFiles := []struct {
		name     string
		path     string
		wantComp bool // si esperamos compresión efectiva
	}{
		{
			name:     "archivo de texto",
			path:     "../testdata/archivo1.txt",
			wantComp: true,
		},
		{
			name:     "archivo markdown",
			path:     "../testdata/archivo2.md",
			wantComp: true,
		},
		{
			name:     "archivo imagen",
			path:     "../testdata/image.jpg",
			wantComp: false, // Las imágenes JPG ya están comprimidas
		},
	}

	// Crear directorio temporal para archivos comprimidos
	tmpDir, err := os.MkdirTemp("", "huffman_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tf := range testFiles {
		t.Run(tf.name, func(t *testing.T) {
			// Leer archivo original
			content, err := os.ReadFile(tf.path)
			if err != nil {
				t.Fatal(err)
			}

			// Comprimir contenido
			encoded, codeMap := Encode(string(content))

			// Guardar versión comprimida
			compressedPath := filepath.Join(tmpDir, filepath.Base(tf.path)+".huf")
			err = os.WriteFile(compressedPath, []byte(encoded), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// Obtener estadísticas
			originalInfo, err := os.Stat(tf.path)
			if err != nil {
				t.Fatal(err)
			}
			compressedInfo, err := os.Stat(compressedPath)
			if err != nil {
				t.Fatal(err)
			}

			// Calcular estadísticas
			originalSize := originalInfo.Size()
			compressedSize := compressedInfo.Size()
			ratio := float64(compressedSize) / float64(originalSize) * 100

			// Mostrar estadísticas detalladas
			t.Logf("Estadísticas para %s:", tf.name)
			t.Logf("  - Archivo: %s", tf.path)
			t.Logf("  - Tamaño original: %d bytes", originalSize)
			t.Logf("  - Tamaño comprimido: %d bytes", compressedSize)
			t.Logf("  - Ratio de compresión: %.2f%%", ratio)
			t.Logf("  - Ahorro de espacio: %d bytes", originalSize-compressedSize)

			// Verificar compresión si se espera
			if tf.wantComp && compressedSize >= originalSize {
				t.Errorf("no se logró la compresión esperada para %s", tf.path)
			}

			// Verificar que podemos recuperar el contenido original
			decoded := Decode(encoded, codeMap)
			if decoded != string(content) {
				t.Error("el contenido decodificado no coincide con el original")
			}
		})
	}
}
