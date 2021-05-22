import type { Todo } from "./api"
import type { AsyncResource } from "./Resource"
import { Text, Small } from "evergreen-ui"
import { RelativeTime } from "./RelativeTime"

import styles from "./TodoList.module.css"

interface TodoListProps {
  resource: AsyncResource<ReadonlyArray<Todo>>
}

function comparator<T, U>(map: (value: T) => U, compare: (a: U, b: U) => number) {
  return function (a: T, b: T) {
    return compare(map(a), map(b))
  }
}

function desc<T>(compare: (a: T, b: T) => number) {
  return function (a: T, b: T) {
    return compare(b, a)
  }
}

export function TodoList(props: TodoListProps) {
  const todos = Array.from(props.resource.read()).sort(
    desc(comparator(_ => _.createdAt, Temporal.Instant.compare))
  )

  return (
    <ul className={styles.todoListContainer}>
      {todos.map(todo => (
        <li key={todo.id} className={styles.todoListItem}>
          <Text>{todo.text}</Text>
          <Text size={300} color="muted">
            <Small>
              <RelativeTime instant={todo.createdAt} />
            </Small>
          </Text>
        </li>
      ))}
    </ul>
  )
}
