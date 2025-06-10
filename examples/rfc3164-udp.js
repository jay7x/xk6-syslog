import syslog from 'k6/x/syslog';

export default function (data) {
    const conn = syslog.connect('localhost:514', {
        transport: 'udp',
    });

    for (let i=0; i<100; i++) {
        const timestamp = new Date().toISOString();
        conn.send(`<134>${timestamp} app-server-k6-vu${__VU} k6: Test log message ${__ITER}\n`);
    }

    conn.close();
}
