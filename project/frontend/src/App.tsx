import React, { useTransition } from "react"
import { createTodo, fetchTodos } from "./api"
import { AsyncResource } from "./Resource"
import { Suspense, useState } from "react"
import { TodoList } from "./TodoList"
import { NewTodoFormlet } from "./NewTodoFormlet"
import { Card, majorScale, Pane } from "evergreen-ui"
import { Spinner, Heading } from "evergreen-ui"

import styles from "./App.module.css"

const initialTodos = new AsyncResource(fetchTodos())

export function App() {
  const [todos, setTodos] = useState(initialTodos)
  const [, startTransition] = useTransition({ timeoutMs: 1000 })

  async function onNewTodoSubmit(text: string) {
    startTransition(async () => {
      await createTodo(text)
      setTodos(new AsyncResource(fetchTodos()))
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
          <Card elevation={0} marginBottom={majorScale(2)}>
            <TodoList resource={todos} />
          </Card>
        </Suspense>
      </div>
    </React.Fragment>
  )
}
