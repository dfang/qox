// https://web.dev/service-worker-mindset/
// https://github.com/GoogleChrome/samples/blob/gh-pages/service-worker/basic/service-worker.js
// https://blog.logrocket.com/every-website-deserves-a-service-worker/
// https://lzw.me/a/pwa-service-worker.html
// https://serviceworke.rs/
// https://medium.com/izettle-engineering/beginners-guide-to-web-push-notifications-using-service-workers-cb3474a17679
// https://love2dev.com/pwa/service-worker-preload/

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
    }).then(() => {
      self.clients.claim();
      // requestNotificationPermission().then(() => {
      //   console.log("request notifition permission")
      // }).catch(() => {
      //   console.log("request failed")
      // })
      // const permission = await requestNotificationPermission();
      // self.Notification.requestPermission();
    })
  )
})

// take care web push notifications events
self.addEventListener('push', function (event) {
  console.log('Received a push message', event);
  const payload = event.data ? event.data.text() : 'no payload';
  event.waitUntil(
    self.registration.showNotification('友情提醒!', {
      body: payload,
    })
  );
});

self.addEventListener('notificationclick', function (event) {
  // Close the notification when it is clicked
  event.notification.close();
  console.log("click event")
  console.log(event)
});

// https://developers.google.com/web/fundamentals/primers/service-workers#cache_and_return_requests
// https://github.com/GoogleChrome/samples/blob/gh-pages/service-worker/read-through-caching/service-worker.js
self.addEventListener('fetch', function (event) {
  event.respondWith(
    caches.match(event.request)
      .then(function (response) {
        // Cache hit - return response
        if (response) {
          return response;
        }
        return fetch(event.request);
      }
      )
  );
});

console.log('Hello from service worker')

const check = () => {
  // if (!('serviceWorker' in navigator)) {
  //   throw new Error('No Service Worker support!')
  // }
  // if (!('PushManager' in window)) {
  //   throw new Error('No Push API Support!')
  // }
}

const requestNotificationPermission = async () => {
  const permission = await window.Notification.requestPermission();
  // value of permission can be 'granted', 'default', 'denied'
  // granted: user has accepted the request
  // default: user has dismissed the notification permission popup by clicking on x
  // denied: user has denied the request.
  if (permission !== 'granted') {
    throw new Error('Permission not granted for Notification');
  }
}

const main = () => {
  check()

  // https://developers.google.com/web/fundamentals/primers/service-workers/registration
  // https://developers.google.com/web/fundamentals/primers/service-workers?hl=zh-CN
  if ('serviceWorker' in navigator) {
    window.addEventListener('load', function () {
      navigator.serviceWorker.register('/admin/assets/javascripts/service-worker.js')
    });
  }
}

main()
