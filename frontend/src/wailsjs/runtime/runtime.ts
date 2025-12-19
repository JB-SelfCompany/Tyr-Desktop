// Wails v2 Runtime API
// Provides access to Wails runtime methods for events, window control, etc.

declare global {
  interface Window {
    runtime: {
      EventsOn(eventName: string, callback: (...args: any[]) => void): () => void;
      EventsEmit(eventName: string, ...args: any[]): void;
      EventsOnce(eventName: string, callback: (...args: any[]) => void): void;
      EventsOff(eventName: string): void;
      EventsOnMultiple(eventName: string, callback: (...args: any[]) => void, maxCallbacks: number): () => void;
      WindowReload(): void;
      WindowSetTitle(title: string): void;
      WindowShow(): void;
      WindowHide(): void;
      WindowMaximise(): void;
      WindowMinimise(): void;
      WindowUnmaximise(): void;
      WindowUnminimise(): void;
      WindowCenter(): void;
      WindowSetSize(width: number, height: number): void;
      WindowGetSize(): Promise<{ width: number; height: number }>;
      WindowSetPosition(x: number, y: number): void;
      WindowGetPosition(): Promise<{ x: number; y: number }>;
      Quit(): void;
      LogPrint(message: string): void;
      LogDebug(message: string): void;
      LogInfo(message: string): void;
      LogWarning(message: string): void;
      LogError(message: string): void;
    };
  }
}

/**
 * EventsOn subscribes to an event
 * @param eventName - The event name to subscribe to
 * @param callback - The callback function to execute when event is received
 * @returns A function to unsubscribe from the event
 */
export function EventsOn(eventName: string, callback: (...args: any[]) => void): () => void {
  if (window.runtime && window.runtime.EventsOn) {
    return window.runtime.EventsOn(eventName, callback);
  }
  console.warn(`EventsOn not available for event: ${eventName}`);
  return () => {};
}

/**
 * EventsEmit emits an event to all listeners
 * @param eventName - The event name to emit
 * @param args - Arguments to pass to the event listeners
 */
export function EventsEmit(eventName: string, ...args: any[]): void {
  if (window.runtime && window.runtime.EventsEmit) {
    window.runtime.EventsEmit(eventName, ...args);
  } else {
    console.warn(`EventsEmit not available for event: ${eventName}`);
  }
}

/**
 * EventsOnce subscribes to an event and unsubscribes after first call
 * @param eventName - The event name to subscribe to
 * @param callback - The callback function to execute once
 */
export function EventsOnce(eventName: string, callback: (...args: any[]) => void): void {
  if (window.runtime && window.runtime.EventsOnce) {
    window.runtime.EventsOnce(eventName, callback);
  } else {
    console.warn(`EventsOnce not available for event: ${eventName}`);
  }
}

/**
 * EventsOff unsubscribes from an event
 * @param eventName - The event name to unsubscribe from
 */
export function EventsOff(eventName: string): void {
  if (window.runtime && window.runtime.EventsOff) {
    window.runtime.EventsOff(eventName);
  } else {
    console.warn(`EventsOff not available for event: ${eventName}`);
  }
}

/**
 * EventsOnMultiple subscribes to an event with a maximum number of callbacks
 * @param eventName - The event name to subscribe to
 * @param callback - The callback function to execute
 * @param maxCallbacks - Maximum number of times callback will be executed
 * @returns A function to unsubscribe from the event
 */
export function EventsOnMultiple(
  eventName: string,
  callback: (...args: any[]) => void,
  maxCallbacks: number
): () => void {
  if (window.runtime && window.runtime.EventsOnMultiple) {
    return window.runtime.EventsOnMultiple(eventName, callback, maxCallbacks);
  }
  console.warn(`EventsOnMultiple not available for event: ${eventName}`);
  return () => {};
}

// Window control functions
export const WindowReload = () => window.runtime?.WindowReload();
export const WindowSetTitle = (title: string) => window.runtime?.WindowSetTitle(title);
export const WindowShow = () => window.runtime?.WindowShow();
export const WindowHide = () => window.runtime?.WindowHide();
export const WindowMaximise = () => window.runtime?.WindowMaximise();
export const WindowMinimise = () => window.runtime?.WindowMinimise();
export const WindowUnmaximise = () => window.runtime?.WindowUnmaximise();
export const WindowUnminimise = () => window.runtime?.WindowUnminimise();
export const WindowCenter = () => window.runtime?.WindowCenter();
export const WindowSetSize = (width: number, height: number) => window.runtime?.WindowSetSize(width, height);
export const WindowGetSize = () => window.runtime?.WindowGetSize();
export const WindowSetPosition = (x: number, y: number) => window.runtime?.WindowSetPosition(x, y);
export const WindowGetPosition = () => window.runtime?.WindowGetPosition();
export const Quit = () => window.runtime?.Quit();

// Logging functions - output to terminal
export const LogPrint = (message: string) => {
  if (window.runtime && window.runtime.LogPrint) {
    window.runtime.LogPrint(message);
  } else {
    console.log(message);
  }
};

export const LogDebug = (message: string) => {
  if (window.runtime && window.runtime.LogDebug) {
    window.runtime.LogDebug(message);
  } else {
    console.debug(message);
  }
};

export const LogInfo = (message: string) => {
  if (window.runtime && window.runtime.LogInfo) {
    window.runtime.LogInfo(message);
  } else {
    console.info(message);
  }
};

export const LogWarning = (message: string) => {
  if (window.runtime && window.runtime.LogWarning) {
    window.runtime.LogWarning(message);
  } else {
    console.warn(message);
  }
};

export const LogError = (message: string) => {
  if (window.runtime && window.runtime.LogError) {
    window.runtime.LogError(message);
  } else {
    console.error(message);
  }
};
