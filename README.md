# Pull-Up Mania

A simple and effective strength training tracker built with Go, HTMX, and SQLite. Perfect for tracking your workouts with minimal complexity and maximum efficiency.

## Features

✅ **Workout Logging**: Easy-to-use interface for logging sets, reps, and weights  
✅ **Admin Panel**: Configure workout plans, exercises, and settings  
✅ **Auto-save**: Progress automatically saves as you type  
✅ **Responsive Design**: Works great on mobile and desktop  
✅ **Local Development**: SQLite for simple local storage  
✅ **Cloud Ready**: Deploy as Google Cloud Function  
✅ **No JavaScript Frameworks**: Pure HTMX for dynamic behavior  

## Quick Start

### Local Development

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd pull-up-mania
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

4. **Open your browser**
   - Visit: http://localhost:8080
   - Admin panel: http://localhost:8080/admin (password: `admin123`)

### First Time Setup

1. **Access Admin Panel**
   - Go to http://localhost:8080/admin
   - Login with password: `admin123`
   - Change the password in settings

2. **Create Workout Plans**
   - Click "Add New Plan"
   - Set the day of week (e.g., Monday = 1)
   - Add exercises with sets and reps
   - Make sure the plan is "Active"

3. **Start Working Out**
   - Go back to home page
   - Log your weights and reps
   - Check off completed sets
   - Add notes as needed

## Project Structure

```
pull-up-mania/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── handlers/        # HTTP handlers (workout, admin)
│   ├── models/          # Data models (Workout, Exercise, Session)
│   ├── storage/         # Storage interface + SQLite implementation
│   └── templates/       # HTML templates with HTMX
├── static/              # CSS and static assets
├── deploy.sh           # Google Cloud deployment script
├── function.yaml       # Cloud Function configuration
└── README.md
```

## Deployment to Google Cloud

### Prerequisites

1. Install [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)
2. Authenticate: `gcloud auth login`
3. Set your project: `gcloud config set project YOUR_PROJECT_ID`

### Deploy

```bash
# Deploy to your Google Cloud project
./deploy.sh YOUR_PROJECT_ID us-central1

# Or use default project
./deploy.sh
```

The script will:
- Deploy your app as a Cloud Function
- Set up HTTP trigger
- Configure memory and timeout settings
- Return the public URL

### Production Notes

- Currently uses SQLite even in production (stored in function instance)
- For production use, implement Firestore storage (code structure ready)
- Consider using Cloud Storage for larger datasets
- Set up proper authentication for admin panel

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)
- `FUNCTION_TARGET`: Set automatically for Cloud Functions
- `GOOGLE_CLOUD_PROJECT`: Used to detect production environment

### Admin Settings

Access via admin panel to configure:
- **Workouts per week**: How many workout days
- **Weight unit**: kg or lbs
- **Admin password**: Change from default

## Technology Stack

- **Backend**: Go with standard `net/http`
- **Frontend**: HTMX + Go templates + Tailwind CSS
- **Storage**: SQLite (Firestore ready for production)
- **Deployment**: Google Cloud Functions

## Why This Stack?

- **Single Language**: Everything in Go, less context switching
- **No Build Step**: Pure server-side rendering, no bundling
- **Lightweight**: Minimal dependencies, fast startup
- **Scalable**: Cloud Functions handle traffic spikes automatically
- **Simple**: Easy to understand and modify

## Development

### Adding New Features

1. **Models**: Add/modify structs in `internal/models/`
2. **Storage**: Update interface in `internal/storage/`
3. **Handlers**: Add routes in `internal/handlers/`
4. **Templates**: Create/modify HTML in `internal/templates/`

### Database Schema

The app uses simple JSON storage in SQLite:
- `settings`: App configuration
- `workout_plans`: Exercise templates for each day
- `workout_sessions`: Individual workout logs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Happy lifting!**