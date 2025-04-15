package manager

//
//// TypeManager manages type-related actions and notifications
//type TypeManager[T any] struct {
//	notifierChan chan string
//	email        string
//	emailToken   string
//	notify       bool
//	prop         c.Property[T]
//	cfg          t.IConfig[T]
//	actions      []t.IAction
//	isRunning    bool
//	mu           sync.Mutex
//}
//
//// NewTypeManager creates a new instance of TypeManager
//func NewTypeManager[T any](cfg t.IConfig[T], propType T) *TypeManager[T] {
//	return &TypeManager[T]{
//		prop:         c.NewProperty[T](reflect.TypeFor[T]().String(), propType),
//		cfg:          cfg,
//		notifierChan: make(chan string, 2),
//		actions:      make([]t.IAction, 0),
//		isRunning:    false,
//	}
//}
//
//// Getters
//func (tm *TypeManager[T]) GetNotifierChan() chan string { return tm.notifierChan }
//func (tm *TypeManager[T]) GetEmail() string             { return tm.email }
//func (tm *TypeManager[T]) GetEmailToken() string        { return tm.emailToken }
//func (tm *TypeManager[T]) GetNotify() bool              { return tm.notify }
//func (tm *TypeManager[T]) GetConfig() t.IConfig[T]      { return tm.cfg }
//func (tm *TypeManager[T]) GetActions() []t.IAction      { return tm.actions }
//func (tm *TypeManager[T]) IsRunning() bool              { return tm.isRunning }
//func (tm *TypeManager[T]) GetProperty() c.Property[T]   { return tm.prop }
//func (tm *TypeManager[T]) GetPropertyType() string      { return reflect.TypeFor[T]().String() }
//func (tm *TypeManager[T]) GetPropertyValue() *T {
//	if v := tm.prop.GetValue(); v != nil {
//		return v
//	} else {
//		return nil
//	}
//}
//func (tm *TypeManager[T]) GetPropertyName() string { return tm.prop.GetName() }
//
//// Setters
//func (tm *TypeManager[T]) SetNotifierChan(notifierChan chan string) { tm.notifierChan = notifierChan }
//func (tm *TypeManager[T]) SetEmail(email string)                    { tm.email = email }
//func (tm *TypeManager[T]) SetEmailToken(emailToken string)          { tm.emailToken = emailToken }
//func (tm *TypeManager[T]) SetNotify(notify bool)                    { tm.notify = notify }
//func (tm *TypeManager[T]) SetConfig(cfg t.IConfig[T])               { tm.cfg = cfg }
//func (tm *TypeManager[T]) AddAction(action t.IAction)               { tm.actions = append(tm.actions, action) }
//
//// StartChecking begins the process of checking Go files
//func (tm *TypeManager[T]) StartChecking(workerCount int) error {
//	if len(tm.actions) == 0 {
//		l.WarnCtx("no actions available to execute", nil)
//		return fmt.Errorf("no actions available to execute")
//	}
//
//	tm.mu.Lock()
//	defer tm.mu.Unlock()
//
//	if tm.isRunning {
//		l.WarnCtx("manager is already running", nil)
//		return fmt.Errorf("manager is already running")
//	}
//
//	l.NoticeCtx("Starting TypeManager", nil)
//	workerManager := NewWorkerManager(workerCount)
//	for _, action := range tm.actions {
//		if action.CanExecute() {
//			l.InfoCtx(fmt.Sprintf("Action %s is executing", action.GetType()), nil)
//			workerManager.GetJobQueue() <- action
//		} else {
//			l.WarnCtx(fmt.Sprintf("Action %s cannot execute", action.GetType()), nil)
//		}
//	}
//
//	go workerManager.StartWorkers()
//	tm.isRunning = true
//	return nil
//}
//func (tm *TypeManager[T]) StopChecking() {
//	//tm.mu.Lock()
//	//defer tm.mu.Unlock()
//
//	if !tm.isRunning {
//		l.WarnCtx("manager is not running", nil)
//		return
//	}
//
//	close(tm.notifierChan)
//	tm.isRunning = false
//	l.InfoCtx("TypeManager stopped successfully", nil)
//}
//func (tm *TypeManager[T]) LoadConfig() error {
//	if tm.cfg == nil {
//		return fmt.Errorf("configuration not initialized")
//	}
//	return nil // Implement loading logic here
//}
//func (tm *TypeManager[T]) SaveConfig() error {
//	if tm.cfg == nil {
//		return fmt.Errorf("configuration not initialized")
//	}
//	return nil // Implement saving logic here
//}
//func (tm *TypeManager[T]) CanNotify() bool {
//	return tm.notify && tm.notifierChan != nil
//}
//func (tm *TypeManager[T]) PrepareActions() error {
//	parsedFiles, err := utils.ParseFiles(tm.cfg.GetPath(), tm.cfg.GetFileType())
//	if err != nil {
//		return fmt.Errorf("error parsing files: %v", err)
//	}
//	// Criar ações baseadas nos arquivos analisados.
//	for pkgName, files := range parsedFiles {
//		action := actions.NewTypeCheckAction(pkgName, files, tm.cfg)
//		tm.AddAction(action)
//	}
//	return nil
//}
