export interface Todo {
  id: string
  text: string
  createdAt: Temporal.Instant
}

interface RawTodo {
  id: string
  text: string
  createdAt: string
}

function deserializeTodo(raw: RawTodo): Todo {
  return {
    id: raw.id,
    text: raw.text,
    createdAt: Temporal.Instant.from(raw.createdAt),
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
  const data = JSON.stringify({ text })
  const response = await fetch("/api/todos", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: data,
  })
  if (!response.ok) {
    throw new Error(response.statusText)
  }
  const json = (await response.json()) as RawTodo
  return deserializeTodo(json)
}
