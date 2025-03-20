<template>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
  <div class="init">
    <!-- Botones en la parte superior -->
    <div class="columns is-centered">
      <!-- Dropdown Archivo -->
      <div class="column is-narrow">
        <div class="dropdown is-hoverable">
          <div class="dropdown-trigger">
            <button class="button" aria-haspopup="true" aria-controls="dropdown-menu">
              <span class="dropdown-text">Archivo</span>
              <span class="icon is-small">
                <i class="fas fa-angle-down" aria-hidden="true"></i>
              </span>
            </button>
          </div>
          <div class="dropdown-menu" id="dropdown-menu" role="menu">
            <div class="dropdown-content">
              <label class="dropdown-item file-input-label"> Abrir Archivo <input type="file" accept=".smia"
                  @change="abrirArchivo" style="display: none;">
              </label>
              <a class="dropdown-item file-input-label" @click="guardar"> Guardar </a>
            </div>
          </div>
        </div>
      </div>

      <!-- Botón Ejecutar -->
      <div class="column is-narrow">
        <button class="button is-primary is-focused" @click="ejecutar">Ejecutar</button>
      </div>

      <!-- Boton Limpiar -->
      <div class="column is-narrow">
        <button class="button is-primary" @click="limpiar">Limpiar</button>
      </div>
    </div>

    <!-- Divs de TextAreas -->
    <div class="columns">
      <!-- Div a la izquierda -->
      <div class="column is-half textarea-wrapper">
        <Codemirror v-model="codigoEntrada" placeholder="ENTRADA DE CODIGO" :extensions="extensions"
          :options="editorOptions"></Codemirror>
      </div>

      <!-- Div a la derecha -->
      <div class="column is-half textarea-wrapper">
        <Codemirror v-model="salidaCodigo" placeholder="SALIDA DE CODIGO" :extensions="outputExtensions"></Codemirror>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { EditorView } from "@codemirror/view";
import { Codemirror } from 'vue-codemirror';
import { EditorState } from "@codemirror/state" // Cambia CodemirrorEditor por Codemirror
export default {
  components: {
    Codemirror // Usamos el componente 'Codemirror'
  },
  data() {
    return {
      codigoEntrada: '',  // Código de entrada desde el textarea
      salidaCodigo: '',   // Resultado desde el backend
      nombreArchivo: '',
      extensions: [
        EditorView.theme({
          "&": {
            backgroundColor: "#282c34",  // Fondo oscuro
            color: "#f8f8f2",         // Color del texto
            height: "550px",
            width: "650px",
            fontSize: "15px",
            textAlign: "left" // Alinear texto a la derecha
          },
          ".cm-gutters": {
            backgroundColor: "#282a36", // Fondo de los números de línea
            color: "#6272a4",            // Color de los números de línea
            border: "none"
          },
          ".cm-cursor": { backgroundColor: "#ffffff" },  // Color del cursor
          ".cm-selection-background": { backgroundColor: "#f1fa8c" }, // Fondo de selección
          ".cm-selection-color": { color: "#282a36" }, // Color de texto seleccionado
          ".cm-activeLine": { backgroundColor: "#44475a" }, // Línea activa
          ".cm-line": { color: "#f8f8f2" }, // Color de texto en línea
        })
      ],
      outputExtensions: [
        EditorState.readOnly.of(true),
        EditorView.theme({
          "&": {
            backgroundColor: "#1e1e1e", // Fondo del output
            color: "#f8f8f2",           // Color del texto
            height: "550px",
            width: "650px",
            fontSize: "17px" ,
            textAlign: "left" // Alinear texto a la derecha
          },
          ".cm-gutters": {
            backgroundColor: "#1e1e1e", // Fondo de los números de línea
            color: "#6272a4",            // Color de los números de línea
            border: "none"
          },
          ".cm-line": { color: "#00ffff" } // Color de texto en el output
        })
      ]
    };
  },
  methods: {

    async ejecutar() {
  try {
    this.salidaCodigo = "";

    // Solicitud POST al backend
    const response = await axios.post('http://localhost:7777/scannear', {
      input: this.codigoEntrada
    });

    this.salidaCodigo = response.data.consola;
  } catch (error) {
    console.error('Error al ejecutar el comando:', error);
    this.salidaCodigo = 'Error en la ejecución del comando';
  }
}

    , abrirArchivo(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.codigoEntrada = e.target.result;
          this.nombreArchivo = file.name; // Actualiza el textarea con el contenido del archivo
        };
        reader.readAsText(file);  // Leer el archivo como texto
      }
    }
    ,
    limpiar() {
      this.salidaCodigo = "";
      this.codigoEntrada = "";
    }
    ,
    guardar() {
      // Si hay un archivo abierto, guardamos los cambios
      if (this.nombreArchivo) {
        const blob = new Blob([this.codigoEntrada], { type: "text/plain" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = this.nombreArchivo; // Usa el nombre del archivo abierto
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
      } else {
        // Si no hay archivo abierto, creamos uno nuevo
      }
    },
  }
};

</script>


<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.init {
  background-color: #1e1e2e;
}

.textarea.is-info.is-medium {
  resize: none;
  width: 550px;
  height: 550px;
  background-color: #282c34;
  color: whitesmoke;
  box-shadow: 4px 4px 10px 2px rgba(23, 179, 236, 0.3);
}

.textarea.is-warning.is-medium {
  resize: none;
  height: 550px;
  background-color: #1e1e1e;
  width: 550px;
  color: #98c379;
  box-shadow: 4px 4px 10px 2px rgba(245, 253, 125, 0.3);
}

.button.is-primary.is-focused {
  background-color: #42c8f9;
  color: #1e1e1e;
  font-style: oblique;
}
.button.is-primary {
  background-color: #1bc40d;
  color: #1e1e1e;
  font-style: oblique;
}

.button {
  background-color: #f4c773;
  color: #1e1e1e;
  font-weight: bold;
}

.dropdown-item {
  background-color: #1a1a1a;
}


.column.is-half.textarea-wrapper {
  display: grid;
  justify-content: center;
  align-items: center;

}

.column.is-narrow {
  margin-top: 40px;
  margin-bottom: 25px;
}

.file-input-label {
  cursor: pointer;
}

.file-input-label2 {
  cursor: pointer;
}
.file-input-label2:hover {
  background-color: #800000;
  /* Color en hover */
}

.file-input-label:hover {
  background-color: #00b89c;
  /* Color en hover */
}

.file-input {
  display: none;
  /* Esconder el input de archivo */
}
</style>
