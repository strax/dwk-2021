enum Status {
  PENDING = "pending",
  ERROR = "error",
  SUCCESS = "success",
}

export class AsyncResource<T> {
  #promise: Promise<void>
  #status: Status = Status.PENDING
  #data?: T
  #error?: Error

  constructor(promise: Promise<T>) {
    this.#promise = promise
      .then(data => {
        this.#status = Status.SUCCESS
        this.#data = data
      })
      .catch(error => {
        this.#status = Status.ERROR
        if (error instanceof Error) {
          this.#error = error
        } else {
          this.#error = new Error(String(error))
        }
      })
  }

  read(): T {
    return this.#readOrSuspend()
  }

  #readOrSuspend() {
    // FIXME: Need to use a temporary to work around a terser bug,
    // see https://github.com/terser/terser/issues/999
    const temp = this.#promise
    switch (this.#status) {
      case Status.PENDING:
        throw temp
      case Status.ERROR:
        throw this.#error
      case Status.SUCCESS:
        return this.#data!
    }
  }
}
