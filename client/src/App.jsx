import { useState, useEffect, useCallback } from "react";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";
import { getTodos, createTodo, updateTodo, deleteTodo } from "./api";
import "./App.css";

function App() {
  const [todos, setTodos] = useState([]);
  const [filter, setFilter] = useState("all");
  const [userId, setUserId] = useState("a2fdc1b2-b6e1-4b7b-97d1-6d1b43f67e1d");
  const [pageSize, setPageSize] = useState(5);
  const [message, setMessage] = useState("");
  const [nextPageToken, setNextPageToken] = useState("");
  const [prevPageToken, setPrevPageToken] = useState("");
  const [pageNo, setPageNo] = useState(1);

  const loadTodos = useCallback(async () => {
    const response = await getTodos({
      user_id: userId,
      status: filter !== "all" ? filter : "",
      page_size: pageSize,
      page_token: nextPageToken,
    });
    console.log("Fetched Todos:", response.todos);
    setTodos(response.todos);
    setNextPageToken(response.next_page_token);
  }, [userId, filter, pageSize]);

  useEffect(() => {
    loadTodos();
  }, [loadTodos]);

  const handleCreate = async (todo) => {
    const newTodo = await createTodo(userId, todo);
    setTodos([...todos, newTodo]);
    setMessage("TODO created successfully.");
    setTimeout(() => setMessage(""), 3000);
  };

  const handleUpdate = async (id, updatedTodo) => {
    await updateTodo(id, updatedTodo);
    setTodos(todos.map((todo) => (todo.id === id ? updatedTodo : todo)));
    setMessage("TODO updated successfully.");
    setTimeout(() => setMessage(""), 3000);
  };

  const handleDelete = async (id) => {
    await deleteTodo(id);
    setTodos(todos.filter((todo) => todo.id !== id));
    setMessage("TODO deleted successfully.");
    setTimeout(() => setMessage(""), 3000);
  };

  const handlePrevPage = async () => {
    const response = await getTodos({
      user_id: userId,
      status: filter !== "all" ? filter : "",
      page_size: pageSize,
      page_token: prevPageToken,
    });
    if (response.todos === null) {
      setTodos([]);
    } else {
      setTodos(response.todos);
      setNextPageToken(response.next_page_token);
      setPrevPageToken(response.prev_page_token);
    }
  };

  const handleNextPage = async () => {
    const response = await getTodos({
      user_id: userId,
      status: filter !== "all" ? filter : "",
      page_size: pageSize,
      page_token: nextPageToken,
    });
    if (response.todos === null) {
      setTodos([]);
    } else {
      setTodos(response.todos);
      setPageNo(pageNo + 1);
      setNextPageToken(response.next_page_token);
      setPrevPageToken(response.prev_page_token);
    }
  };

  const handlePageSizeChange = (event) => {
    setPageSize(Number(event.target.value));
    setNextPageToken("");
    setPrevPageToken("");
  };

  return (
    <div className="App">
      <h1>TODO App</h1>
      {message && <p className="message">{message}</p>}
      <TodoForm onCreate={handleCreate} />
      <div>
        <button onClick={() => setFilter("all")}>All</button>
        <button onClick={() => setFilter("pending")}>Pending</button>
        <button onClick={() => setFilter("completed")}>Completed</button>
      </div>
      {todos.length > 0 ? (
        <TodoList
          todos={todos}
          onUpdate={handleUpdate}
          onDelete={handleDelete}
        />
      ) : (
        <p>No TODOs available</p>
      )}
      <div>
        <button onClick={handlePrevPage} disabled={pageNo === 1}>
          Previous
        </button>
        <button onClick={handleNextPage} disabled={!nextPageToken}>
          Next
        </button>
        Todo per page :{" "}
        <select value={pageSize} onChange={handlePageSizeChange}>
          <option value={5}>5</option>
          <option value={10}>10</option>
          <option value={20}>20</option>
        </select>
      </div>
    </div>
  );
}

export default App;
