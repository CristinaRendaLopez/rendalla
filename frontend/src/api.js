import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'https://ji3pj4w8rg.execute-api.eu-north-1.amazonaws.com',
  headers: {
    'Content-Type': 'application/json',
  },
});

export default apiClient;
