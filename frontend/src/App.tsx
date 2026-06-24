import './App.css'

function App() {
  return (
    <div className="app">
      <header className="header">
        <h1>time-broker</h1>
        <p className="subtitle">Get back your time. Use a managed broker, hooked up to your calendar provider.</p>
      </header>

      <main className="main">
        <section className="card">
          <h2>Calendar</h2>
          <p className="placeholder">Your events will appear here once connected to a calendar provider.</p>
        </section>

        <section className="card">
          <h2>Providers</h2>
          <ul className="provider-list">
            <li>Google Calendar</li>
          </ul>
        </section>
      </main>
    </div>
  )
}

export default App
