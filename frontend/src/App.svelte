<script lang="ts">
  import { onMount } from 'svelte'

  type Todo = {
    id: number
    title: string
    completed: boolean
    created_at: string
  }

  let todos: Todo[] = []
  let newTitle = ''
  let error = ''
  const API_URL = 'http://localhost:8080/todos'

  async function loadTodos() {
    try {
      const response = await fetch(API_URL)
      if (!response.ok) throw new Error('Failed to load todos')
      todos = await response.json()
    } catch (err) {
      error = String(err)
    }
  }

  async function addTodo() {
    if (!newTitle.trim()) return

    const response = await fetch(API_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle.trim() }),
    })

    if (response.ok) {
      newTitle = ''
      await loadTodos()
      error = ''
    } else {
      const body = await response.json()
      error = body.error || 'Could not add todo'
    }
  }

  async function toggleTodo(todo: Todo) {
    const response = await fetch(`${API_URL}/${todo.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ completed: !todo.completed }),
    })
    if (response.ok) {
      await loadTodos()
      error = ''
    } else {
      const body = await response.json()
      error = body.error || 'Could not update todo'
    }
  }

  async function deleteTodo(id: number) {
    const response = await fetch(`${API_URL}/${id}`, {
      method: 'DELETE',
    })
    if (response.ok) {
      todos = todos.filter((todo) => todo.id !== id)
      error = ''
    } else {
      const body = await response.json()
      error = body.error || 'Could not remove todo'
    }
  }

  onMount(loadTodos)
</script>

<main>
  <section class="todo-app">
    <h1>Todo List</h1>

    <form on:submit|preventDefault={addTodo} class="todo-form">
      <input
        type="text"
        bind:value={newTitle}
        placeholder="Add a new todo"
        aria-label="New todo title"
      />
      <button type="submit">Add</button>
    </form>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <ul class="todo-list">
      {#each todos as todo}
        <li class:completed={todo.completed}>
          <label>
            <input
              type="checkbox"
              checked={todo.completed}
              on:change={() => toggleTodo(todo)}
            />
            <span>{todo.title}</span>
          </label>
          <button class="delete" on:click={() => deleteTodo(todo.id)}>×</button>
        </li>
      {/each}
    </ul>
  </section>
</main>

<style>
  main {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    background: #f8f9fb;
    color: #111827;
  }

  .todo-app {
    width: min(680px, 100%);
    background: #ffffff;
    border-radius: 1rem;
    padding: 2rem;
    box-shadow: 0 24px 60px rgba(15, 23, 42, 0.08);
  }

  h1 {
    margin: 0 0 1rem;
    font-size: 2rem;
  }

  .todo-form {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  input[type='text'] {
    flex: 1;
    padding: 0.85rem 1rem;
    border: 1px solid #d1d5db;
    border-radius: 0.75rem;
    font-size: 1rem;
  }

  button {
    border: none;
    background: #2563eb;
    color: white;
    padding: 0.85rem 1.25rem;
    border-radius: 0.75rem;
    cursor: pointer;
    font-weight: 600;
  }

  .error {
    margin-bottom: 1rem;
    color: #b91c1c;
  }

  .todo-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.75rem;
  }

  .todo-list li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.95rem 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.85rem;
    background: #f8fafc;
  }

  .todo-list li.completed span {
    text-decoration: line-through;
    opacity: 0.65;
  }

  .todo-list label {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
  }

  .delete {
    border: none;
    background: transparent;
    color: #ef4444;
    font-size: 1.25rem;
    cursor: pointer;
  }
</style>
