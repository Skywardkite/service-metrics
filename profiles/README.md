До:
      flat  flat%   sum%        cum   cum%
    2565kB 49.67% 49.67%     2565kB 49.67%  runtime.allocm
 1184.27kB 22.93% 72.61%  1184.27kB 22.93%  runtime/pprof.StartCPUProfile
  902.59kB 17.48% 90.08%   902.59kB 17.48%  compress/flate.NewWriter (inline)
  512.05kB  9.92%   100%   512.05kB  9.92%  time.NewTicker
         0     0%   100%   902.59kB 17.48%  compress/gzip.(*Writer).Write

Проблема, которую правила:
    2565kB 49.67% 49.67%     2565kB 49.67%  runtime.allocm
         0     0%   100%   902.59kB 17.48%  compress/gzip.(*Writer).Write

В func (w *gzipResponseWriter) Write -> WriteHeader я создавала gzip.NewWriter(w.ResponseWriter) 
получается каждый раз аллоцировали память для создания еще одного gzip.Writer при вызове Write, что избыточно

Что делала:
добавила gzipPool = sync.Pool. Тпереь не будет происходить постоянное создание нового writer, будем использовать из пулла. 
Если в пулле нет writer, то он будет создан. Это меньше аллокаций чем раньше, поэтом при сравнении видим 
    -513kB  9.93% 30.44%     -513kB  9.93%  runtime.allocm
         0     0% 20.52%  -902.59kB 17.48%  compress/gzip.(*Writer).Write


File: server
Type: inuse_space
Time: 2025-12-16 23:11:29 MSK
Showing nodes accounting for -1059.55kB, 20.52% of 5163.92kB total
      flat  flat%   sum%        cum   cum%
-1184.27kB 22.93% 22.93% -1184.27kB 22.93%  runtime/pprof.StartCPUProfile
    1028kB 19.91%  3.03%     1028kB 19.91%  bufio.NewReaderSize (inline)
 -902.59kB 17.48% 20.50%  -902.59kB 17.48%  compress/flate.NewWriter (inline)
    -513kB  9.93% 30.44%     -513kB  9.93%  runtime.allocm
  512.22kB  9.92% 20.52%   512.22kB  9.92%  runtime.malg
  512.14kB  9.92% 10.60%   512.14kB  9.92%  github.com/golang-migrate/migrate/v4/source.Register
 -512.05kB  9.92% 20.52%  -512.05kB  9.92%  time.NewTicker
         0     0% 20.52%     1028kB 19.91%  bufio.NewReader (inline)
         0     0% 20.52%  -902.59kB 17.48%  compress/gzip.(*Writer).Write
         0     0% 20.52%  -512.05kB  9.92%  github.com/Skywardkite/service-metrics/internal/filestorage.(*StorageConfig).Run.func1
         0     0% 20.52%  -902.59kB 17.48%  github.com/Skywardkite/service-metrics/internal/handler.(*Handler).UpdateJSONHandler
         0     0% 20.52% -2086.86kB 40.41%  github.com/Skywardkite/service-metrics/internal/logger.WithLogging.func1
         0     0% 20.52% -1184.27kB 22.93%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0% 20.52% -2086.86kB 40.41%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 20.52% -2086.86kB 40.41%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 20.52%   512.14kB  9.92%  github.com/golang-migrate/migrate/v4/source/file.init.0
         0     0% 20.52%  -902.59kB 17.48%  main.(*gzipResponseWriter).Write
         0     0% 20.52% -2086.86kB 40.41%  main.gzipMiddleware.func1
         0     0% 20.52% -2086.86kB 40.41%  main.main.func2.authMiddleware.1
         0     0% 20.52% -1058.85kB 20.50%  net/http.(*conn).serve
         0     0% 20.52% -2086.86kB 40.41%  net/http.HandlerFunc.ServeHTTP
         0     0% 20.52%     1028kB 19.91%  net/http.newBufioReader
         0     0% 20.52% -2086.86kB 40.41%  net/http.serverHandler.ServeHTTP
         0     0% 20.52% -1184.27kB 22.93%  net/http/pprof.Profile
         0     0% 20.52%   512.14kB  9.92%  runtime.doInit (inline)
         0     0% 20.52%   512.14kB  9.92%  runtime.doInit1
         0     0% 20.52%      513kB  9.93%  runtime.handoffp
         0     0% 20.52%   512.14kB  9.92%  runtime.main
         0     0% 20.52%    -1026kB 19.87%  runtime.mcall
         0     0% 20.52%      513kB  9.93%  runtime.mstart
         0     0% 20.52%      513kB  9.93%  runtime.mstart0
         0     0% 20.52%      513kB  9.93%  runtime.mstart1
         0     0% 20.52%     -513kB  9.93%  runtime.newm
         0     0% 20.52%   512.22kB  9.92%  runtime.newproc.func1
         0     0% 20.52%   512.22kB  9.92%  runtime.newproc1
         0     0% 20.52%    -1026kB 19.87%  runtime.park_m
         0     0% 20.52%    -1026kB 19.87%  runtime.resetspinning
         0     0% 20.52%      513kB  9.93%  runtime.retake
         0     0% 20.52%    -1026kB 19.87%  runtime.schedule
         0     0% 20.52%     -513kB  9.93%  runtime.startm
         0     0% 20.52%      513kB  9.93%  runtime.sysmon
         0     0% 20.52%   512.22kB  9.92%  runtime.systemstack
         0     0% 20.52%    -1026kB 19.87%  runtime.wakep