import { Component, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error("Enque UI Error:", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex flex-col items-center justify-center h-screen p-8" style={{ background: '#0a0a0f', color: '#e8e6e3' }}>
          <div className="w-12 h-12 rounded-xl flex items-center justify-center mb-5" style={{ background: 'rgba(248, 113, 113, 0.1)', border: '1px solid rgba(248, 113, 113, 0.2)' }}>
            <span className="text-xl" style={{ color: '#f87171' }}>!</span>
          </div>
          <h1 className="text-base font-display font-bold mb-4" style={{ color: '#f87171' }}>
            An unexpected error occurred
          </h1>
          <pre
            className="text-xs font-mono p-4 rounded-lg max-w-lg overflow-auto mb-6"
            style={{
              background: 'rgba(16, 16, 22, 0.8)',
              color: '#9d9da7',
              border: '1px solid rgba(255,255,255,0.06)',
            }}
          >
            {this.state.error?.message}
          </pre>
          <button
            onClick={() => this.setState({ hasError: false, error: null })}
            className="btn-primary"
          >
            Try Again
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
