export interface Transition<Ctx, State extends string, E> {
  target: State
  /** Naya context banata — di gayi to update, warna context waisa hi rehta. */
  assign?: (context: Ctx, event: E) => Ctx
}

type TransitionConfig<Ctx, State extends string, E> = State | Transition<Ctx, State, E>

interface StateNode<Ctx, State extends string, Event extends { type: string }> {
  on?: {
    [K in Event['type']]?: TransitionConfig<Ctx, State, Extract<Event, { type: K }>>
  }
}

export interface MachineConfig<Ctx, State extends string, Event extends { type: string }> {
  initial: State
  context: Ctx
  states: Record<State, StateNode<Ctx, State, Event>>
}

export interface MachineSnapshot<Ctx, State extends string> {
  state: State
  context: Ctx
}

export interface Machine<Ctx, State extends string, Event extends { type: string }> {
  initial: MachineSnapshot<Ctx, State>
  /** Ek event apply karta — jo state se woh event defined nahi, snapshot waisa hi laut aata hai. */
  transition: (snapshot: MachineSnapshot<Ctx, State>, event: Event) => MachineSnapshot<Ctx, State>
  /** Is state se ye event handle hota hai kya — Button disable karne jaisi checks ke liye. */
  can: (state: State, eventType: Event['type']) => boolean
}

/**
 * Ek chhota, dependency-free finite state machine. Reach for it jab kayi
 * boolean flags (isLoading/isOpen/hasError/...) combine ho ke "impossible
 * states" bana dete hain — explicit states + transitions un combinations ko
 * hi allow karte hain jo actually valid hain.
 *
 * `states`/`context` sirf plain data hai — koi React yahan nahi. Component
 * ke andar `useMachine` (`src/hooks/useMachine.ts`) se jodo.
 *
 * @example
 * type State = 'idle' | 'loading' | 'success' | 'error'
 * type Event =
 *   | { type: 'FETCH' }
 *   | { type: 'RESOLVE'; data: string }
 *   | { type: 'REJECT'; message: string }
 *   | { type: 'RETRY' }
 *
 * const fetchMachine = createMachine<{ data?: string; error?: string }, State, Event>({
 *   initial: 'idle',
 *   context: {},
 *   states: {
 *     idle: { on: { FETCH: 'loading' } },
 *     loading: {
 *       on: {
 *         RESOLVE: { target: 'success', assign: (ctx, e) => ({ ...ctx, data: e.data }) },
 *         REJECT: { target: 'error', assign: (ctx, e) => ({ ...ctx, error: e.message }) },
 *       },
 *     },
 *     success: { on: { FETCH: 'loading' } },
 *     error: { on: { RETRY: 'loading' } },
 *   },
 * })
 */
export function createMachine<Ctx, State extends string, Event extends { type: string }>(
  config: MachineConfig<Ctx, State, Event>,
): Machine<Ctx, State, Event> {
  function transition(
    snapshot: MachineSnapshot<Ctx, State>,
    event: Event,
  ): MachineSnapshot<Ctx, State> {
    const node = config.states[snapshot.state]
    const raw = node.on?.[event.type as Event['type']]
    if (!raw) return snapshot

    // `on[event.type]` ki exact-event narrowing yahan static rahke type-check
    // nahi ho sakti (dynamic key lookup) — runtime par event.type se hi milaya hai.
    const entry = raw as TransitionConfig<Ctx, State, Event>
    if (typeof entry === 'string') {
      return { state: entry, context: snapshot.context }
    }
    return {
      state: entry.target,
      context: entry.assign ? entry.assign(snapshot.context, event) : snapshot.context,
    }
  }

  function can(state: State, eventType: Event['type']): boolean {
    return Boolean(config.states[state].on?.[eventType])
  }

  return { initial: { state: config.initial, context: config.context }, transition, can }
}
