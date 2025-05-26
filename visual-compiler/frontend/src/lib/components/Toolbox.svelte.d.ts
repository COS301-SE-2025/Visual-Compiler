import { SvelteComponentTyped } from 'svelte';

export default class Toolbox extends SvelteComponentTyped<
  // no props
  {},
  // these are the events: blockClick carries a string
  { blockClick: string },
  // no slots
  {}
> {}
