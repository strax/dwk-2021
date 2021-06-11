import type { Todo } from "./api"
import type { AsyncResource } from "./Resource"
import { Text, Small, Checkbox, Card, majorScale } from "evergreen-ui"
import { RelativeTime } from "./RelativeTime"

import styles from "./TodoList.module.css"

interface TodoListProps {
  resource: AsyncResource<ReadonlyArray<Todo>>
  onDoneChange: (id: string, done: boolean) => void
}

type Comparator<T> = (a: T, b: T) => number

namespace Comparator {
  export function map<T, U>(map: (value: T) => U, compare: Comparator<U>): Comparator<T> {
    return function (a: T, b: T) {
      return compare(map(a), map(b))
    }
  }

  export function desc<T>(compare: Comparator<T>): Comparator<T> {
    return function (a: T, b: T) {
      return compare(b, a)
    }
  }
}

function partition<T>(xs: Iterable<T>, predicate: (x: T) => boolean): readonly [T[], T[]] {
  const p1 = []
  const p2 = []
  for (const x of xs) {
    if (predicate(x)) {
      p1.push(x)
    } else {
      p2.push(x)
    }
  }
  return [p1, p2] as const
}

interface TodoListItemProps {
  resource: Todo
  onDoneChange: (done: boolean) => void
}

function TodoListItem({ resource, onDoneChange }: TodoListItemProps) {
  return (
    <li className={styles.todoListItem}>
      <Checkbox checked={resource.done} onChange={evt => onDoneChange(evt.currentTarget.checked)} />
      <div className={styles.todoListItemText}>
        {resource.done ? (
          <Text textDecoration="line-through" color="muted">
            {resource.text}
          </Text>
        ) : (
          <Text>{resource.text}</Text>
        )}

        <Text size={300} color="muted">
          <Small>
            <RelativeTime instant={resource.createdAt} />
          </Small>
        </Text>
      </div>
    </li>
  )
}

export function TodoList(props: TodoListProps) {
  const [doneTodos, notDoneTodos] = partition(props.resource.read(), todo => todo.done)
  doneTodos.sort(Comparator.desc(Comparator.map(_ => _.updatedAt, Temporal.Instant.compare)))
  notDoneTodos.sort(Comparator.desc(Comparator.map(_ => _.createdAt, Temporal.Instant.compare)))

  return (
    <>
      <Card elevation={0} marginBottom={majorScale(2)}>
        <ul className={styles.todoListContainer}>
          {notDoneTodos.map(todo => (
            <TodoListItem
              key={todo.id}
              resource={todo}
              onDoneChange={done => props.onDoneChange(todo.id, done)}
            />
          ))}
        </ul>
      </Card>
      <Card elevation={0} marginBottom={majorScale(2)}>
        <ul className={styles.todoListContainer}>
          {doneTodos.map(todo => (
            <TodoListItem
              key={todo.id}
              resource={todo}
              onDoneChange={done => props.onDoneChange(todo.id, done)}
            />
          ))}
        </ul>
      </Card>
    </>
  )
}
