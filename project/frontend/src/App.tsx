import React, { useTransition, Suspense, useState } from "react"
import { createTodo, fetchTodos, updateTodo } from "./api"
import { AsyncResource } from "./Resource"
import { TodoList } from "./TodoList"
import { NewTodoFormlet } from "./NewTodoFormlet"
import { majorScale, Pane, Spinner, Heading, toaster } from "evergreen-ui"

import styles from "./App.module.css"

const initialTodos = new AsyncResource(fetchTodos())

function showErrorToast(error: unknown) {
  console.error(error)
  toaster.danger("Something went wrong", {
    description: error instanceof Error ? error.message : String(error),
  })
}

export function App() {
  const [todos, setTodos] = useState(initialTodos)
  const [, startTransition] = useTransition({ timeoutMs: 1000 })

  async function onNewTodoSubmit(text: string) {
    startTransition(async () => {
      await createTodo(text)
      setTodos(new AsyncResource(fetchTodos()))
    })
  }

  async function onTodoDoneChange(id: string, done: boolean) {
    startTransition(async () => {
      try {
        await updateTodo(id, { done })
        setTodos(new AsyncResource(fetchTodos()))
      } catch (err) {
        showErrorToast(err)
      }
    })
  }

  return (
    <React.Fragment>
      <div className={styles.container}>
        <img className={styles.heroImage} src="/api/image" width={100} height={100} />
        <Heading is="h1" textAlign="center" size={900} marginY={majorScale(2)} userSelect="none">
          Todos
        </Heading>
        <Pane paddingBottom={majorScale(2)} marginBottom={1} background="white">
          <NewTodoFormlet onSubmit={onNewTodoSubmit} />
        </Pane>
        <Suspense
          fallback={
            <Pane display="flex" alignItems="center" justifyContent="center" height={100}>
              <Spinner />
            </Pane>
          }>
          <TodoList resource={todos} onDoneChange={onTodoDoneChange} />
        </Suspense>
      </div>
    </React.Fragment>
  )
}
