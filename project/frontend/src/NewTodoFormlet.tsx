import React, { useState } from "react"
import { majorScale, TextInput } from "evergreen-ui"

import styles from "./NewTodoFormlet.module.css"

interface NewTodoFormletProps {
  onSubmit: (text: string) => Promise<void> | void
}

export function NewTodoFormlet(props: NewTodoFormletProps) {
  const [draftText, setDraftText] = useState("")
  const [inFlight, setInFlight] = useState(false)
  const disabled = inFlight || draftText.trim() === ""

  async function submit(evt: React.FormEvent<HTMLFormElement>) {
    evt.preventDefault()
    if (disabled) return

    setInFlight(true)
    try {
      await props.onSubmit(draftText)
      setDraftText("")
    } finally {
      setInFlight(false)
    }
  }

  return (
    <form onSubmit={submit}>
      <TextInput
        height={majorScale(5)}
        width="100%"
        autoFocus
        placeholder="Type and press Enter"
        onChange={(evt: React.ChangeEvent<HTMLInputElement>) => setDraftText(evt.target.value)}
        value={draftText}
        maxLength={140}
      />
      <button className={styles.submitButton} type="submit" disabled={disabled}>
        Create
      </button>
    </form>
  )
}
