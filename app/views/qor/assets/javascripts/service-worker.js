// https://github.com/GoogleChrome/samples/blob/gh-pages/service-worker/basic/service-worker.js
// https://lzw.me/a/pwa-service-worker.html
// https://serviceworke.rs/
// https://medium.com/izettle-engineering/beginners-guide-to-web-push-notifications-using-service-workers-cb3474a17679



const PRECACHE = 'precache-v1';
const RUNTIME = 'runtime';

// A list of local resources we always want to be cached.
const PRECACHE_URLS = [
  '/admin/assets/images/logo.png',
  '/favicon.ico',
  '/admin/assets/javascripts/vendors.js',
  '/admin/assets/javascripts/qor_admin_default.js',
  '/admin/assets/stylesheets/qor_admin_default.css',
  '/admin/assets/javascripts/qor_demo.js',
  '/admin/assets/stylesheets/qor_demo.css',
  '/admin/assets/javascripts/notifications.js',
  '/admin/assets/stylesheets/notifications.css',
  '/admin/assets/fonts/MaterialIcons-Regular.woff2',
  '/admin/assets/fonts/MaterialIcons-Regular.woff',
  '/admin/assets/fonts/MaterialIcons-Regular.ttf',
  '/admin/assets/fonts/Roboto-BoldItalic.ttf',
  '/admin/assets/fonts/Roboto-Bold.ttf',
  '/admin/assets/fonts/Roboto-MediumItalic.ttf',
  '/admin/assets/fonts/Roboto-Medium.ttf',
  '/admin/assets/fonts/Roboto-Italic.ttf',
  '/admin/assets/fonts/Roboto-Regular.ttf',
  '/admin/assets/fonts/Roboto-Light.ttf',
  '/admin/assets/fonts/Roboto-LightItalic.ttf',
  '/admin/assets/fonts/Roboto-Thin.ttf',
  '/admin/assets/fonts/Roboto-ThinItalic.ttf',
];

// The install handler takes care of precaching the resources we always need.
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(PRECACHE)
      .then(cache => cache.addAll(PRECACHE_URLS))
      .then(self.skipWaiting())
  );
});

// The activate handler takes care of cleaning up old caches.
self.addEventListener('activate', event => {
  const currentCaches = [PRECACHE, RUNTIME];
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return cacheNames.filter(cacheName => !currentCaches.includes(cacheName));
    }).then(cachesToDelete => {
      return Promise.all(cachesToDelete.map(cacheToDelete => {
        return caches.delete(cacheToDelete);
      }));
    }).then(() => self.clients.claim())
  );
});

