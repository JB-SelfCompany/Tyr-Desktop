/**
 * useEventStream - Hook for subscribing to Wails runtime events
 * Handles real-time updates from backend via EventsOn
 */

import { useEffect } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { useServiceStore } from '../store/serviceStore';
import { useLogsStore } from '../store/logsStore';
import type {
  LogEventDTO,
  MailEventDTO,
  ConnectionEventDTO,
  ServiceStatusDTO,
  PeerInfoDTO,
} from '../wailsjs/go/main/models';

/**
 * Event names emitted by the backend
 */
export const EventNames = {
  SERVICE_LOG: 'service:log',
  SERVICE_MAIL: 'service:mail',
  SERVICE_CONNECTION: 'service:connection',
  SERVICE_STATUS: 'service:status',
  SERVICE_PEERS: 'service:peers',
} as const;

/**
 * Hook that subscribes to service log events
 * Automatically adds logs to the logs store
 */
export function useLogEvents() {
  const addLog = useLogsStore((state) => state.addLog);

  useEffect(() => {
    const unsubscribe = EventsOn(EventNames.SERVICE_LOG, (log: LogEventDTO) => {
      addLog(log);
    });

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [addLog]);
}

/**
 * Hook that subscribes to mail events
 * Can be used to show notifications for new mail
 */
export function useMailEvents(
  onMail?: (mail: MailEventDTO) => void
) {
  useEffect(() => {
    if (!onMail) return;

    const unsubscribe = EventsOn(EventNames.SERVICE_MAIL, (mail: MailEventDTO) => {
      onMail(mail);
    });

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [onMail]);
}

/**
 * Hook that subscribes to connection events
 * Can be used to show notifications for peer connections
 */
export function useConnectionEvents(
  onConnection?: (connection: ConnectionEventDTO) => void
) {
  useEffect(() => {
    if (!onConnection) return;

    const unsubscribe = EventsOn(
      EventNames.SERVICE_CONNECTION,
      (connection: ConnectionEventDTO) => {
        onConnection(connection);
      }
    );

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [onConnection]);
}

/**
 * Hook that subscribes to service status updates
 * Automatically updates the service store
 */
export function useServiceStatusEvents() {
  const setStatus = useServiceStore((state) => state.setStatus);

  useEffect(() => {
    const unsubscribe = EventsOn(
      EventNames.SERVICE_STATUS,
      (status: ServiceStatusDTO) => {
        setStatus(status);
      }
    );

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [setStatus]);
}

/**
 * Hook that subscribes to peer statistics updates
 * Automatically updates the service store
 */
export function usePeerStatsEvents() {
  const setPeers = useServiceStore((state) => state.setPeers);

  useEffect(() => {
    const unsubscribe = EventsOn(
      EventNames.SERVICE_PEERS,
      (peers: PeerInfoDTO[]) => {
        setPeers(peers);
      }
    );

    return () => {
      if (unsubscribe) unsubscribe();
    };
  }, [setPeers]);
}

/**
 * Master hook that subscribes to all event streams
 * Use this in your root App component to enable all real-time updates
 *
 * @example
 * ```tsx
 * function App() {
 *   useEventStream();
 *   return <YourApp />;
 * }
 * ```
 */
export function useEventStream() {
  useLogEvents();
  useServiceStatusEvents();
  usePeerStatsEvents();
}

/**
 * Hook that subscribes to all events with custom handlers
 * Useful when you need custom behavior for each event type
 *
 * @example
 * ```tsx
 * useEventStreamWithHandlers({
 *   onMail: (mail) => {
 *     showNotification('New mail', mail.subject);
 *   },
 *   onConnection: (conn) => {
 *     if (conn.type === 'connected') {
 *       showSuccess('Peer connected', conn.peer);
 *     }
 *   },
 * });
 * ```
 */
export function useEventStreamWithHandlers(handlers: {
  onMail?: (mail: MailEventDTO) => void;
  onConnection?: (connection: ConnectionEventDTO) => void;
}) {
  useEventStream(); // Subscribe to default handlers
  useMailEvents(handlers.onMail);
  useConnectionEvents(handlers.onConnection);
}
