import React from "react";
import "./TodoItem.css";

function TodoItem({ todo, onUpdate, onDelete }) {
  const handleStatusChange = (e) => {
    onUpdate(todo.id, { ...todo, status: e.target.value });
  };

  return (
    <li>
      <div className="Todo-content">
        <h3>{todo.title}</h3>
        <p>{todo.description}</p>
      </div>
      <div className="TodoButton">
        <select value={todo.status} onChange={handleStatusChange}>
          <option value="pending">Pending</option>
          <option value="completed">Completed</option>
        </select>
        <button onClick={() => onDelete(todo.id)}>Delete</button>
      </div>
    </li>
  );
}

export default TodoItem;
