<template>
  <div>
    <h1>Bienvenid@ a Rendalla</h1>
    <h2>Lista de Canciones:</h2>
    <ul>
      <li v-for="song in songs" :key="song.id">
        {{ song.title }} - {{ song.author }}
      </li>
    </ul>
  </div>
</template>

<script>
import apiClient from './api';

export default {
  name: "App",
  data() {
    return {
      songs: [],
    };
  },
  methods: {
    async fetchSongs() {
      try {
        const response = await apiClient.get('/dev/songs');
        this.songs = response.data.data;
        console.log(this.songs);
      } catch (error) {
        console.error('Error al obtener canciones:', error);
      }
    },
  },
  created() {
    this.fetchSongs();
  },
};
</script>

<style>
/* Estilos básicos (opcional) */
h2 {
  margin-top: 20px;
  font-family: Arial, sans-serif;
}
</style>
