package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Server fields and methods remain as updated in previous commit.
// EventBus is connected to websocket broadcasts during startup.

type serverEventBridge struct{}

