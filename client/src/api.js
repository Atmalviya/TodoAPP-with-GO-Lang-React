// const dotenv = require('dotenv').config();
const API_URL = 'http://localhost:8000';

export const getTodos = async (params) => {
  const url = new URL(`${API_URL}/todos`);
  Object.keys(params).forEach(key => url.searchParams.append(key, params[key]));
  const response = await fetch(url);
  
  return response.json();
};

export const createTodo = async (userId, todo) => {
    const response = await fetch(`${API_URL}/todos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...todo, user_id: userId })
    });
    return response.json();
  };

export const updateTodo = async (id, todo) => {
  await fetch(`${API_URL}/todos/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(todo)
  });
};

export const deleteTodo = async (id) => {
  await fetch(`${API_URL}/todos/${id}`, { method: 'DELETE' });
};
