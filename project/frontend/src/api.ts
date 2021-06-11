export interface Todo {
  id: string
  text: string
  done: boolean
  createdAt: Temporal.Instant,
  updatedAt: Temporal.Instant
}

interface RawTodo {
  id: string
  text: string
  createdAt: string
  updatedAt: string
  done: boolean
}

function deserializeTodo(raw: RawTodo): Todo {
  return {
    id: raw.id,
    text: raw.text,
    done: raw.done,
    createdAt: Temporal.Instant.from(raw.createdAt),
    updatedAt: Temporal.Instant.from(raw.updatedAt)
  }
}

function delay(timeout: number): Promise<void> {
  return new Promise(resolve => setTimeout(() => resolve(), timeout))
}

export async function fetchTodos(): Promise<ReadonlyArray<Todo>> {
  const response = await fetch("/api/todos")
  const todos = (await response.json()) as RawTodo[]
  // await delay(2000)
  return todos.map(deserializeTodo)
}

export async function createTodo(text: string): Promise<Todo> {
  const response = await fetch("/api/todos", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) {
    throw new Error(response.statusText)
  }
  const json = (await response.json()) as RawTodo
  return deserializeTodo(json)
}

interface UpdateTodoData {
  done: boolean
}

export async function updateTodo(id: string, data: UpdateTodoData): Promise<Todo> {
  const response = await fetch(`/api/todos/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(data)
  })
  if (!response.ok) {
    throw new Error(response.statusText)
  }
  const json = (await response.json()) as RawTodo
  return deserializeTodo(json)
}
