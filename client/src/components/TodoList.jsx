import React from 'react';
import TodoItem from './TodoItem';
import './TodoList.css';

function TodoList({ todos, onUpdate, onDelete }) {
  console.log('TodoList props:', todos); // Log the todos received
  return (
    <ul>
      {todos.map(todo => (
        <TodoItem key={todo.id} todo={todo} onUpdate={onUpdate} onDelete={onDelete} />
      ))}
    </ul>
  );
}

export default TodoList;
