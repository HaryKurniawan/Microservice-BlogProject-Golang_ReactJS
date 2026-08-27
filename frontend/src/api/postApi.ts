import axios from 'axios';
import type { Post, CreatePostPayload, UpdatePostPayload } from '../types/post';

const BASE_URL = 'http://localhost:8085/api';

const api = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
});

// Menyisipkan token ke setiap request jika ada
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const getAllPosts = async (): Promise<Post[]> => {
  const res = await api.get<Post[]>('/posts');
  return res.data;
};

export const getPostById = async (id: number): Promise<Post> => {
  const res = await api.get<Post>(`/posts/${id}`);
  return res.data;
};

export const createPost = async (payload: CreatePostPayload): Promise<Post> => {
  const res = await api.post<Post>('/posts', payload);
  return res.data;
};

export const updatePost = async (id: number, payload: UpdatePostPayload): Promise<Post> => {
  const res = await api.put<Post>(`/posts/${id}`, payload);
  return res.data;
};

export const deletePost = async (id: number): Promise<void> => {
  await api.delete(`/posts/${id}`);
};

// API calls untuk User Service (Autentikasi)
export const registerUser = async (payload: any): Promise<any> => {
  const res = await api.post('/users/register', payload);
  return res.data;
};

export const loginUser = async (payload: any): Promise<any> => {
  const res = await api.post('/users/login', payload);
  return res.data;
};

export const getUserProfile = async (): Promise<any> => {
  const res = await api.get('/users/profile');
  return res.data;
};
