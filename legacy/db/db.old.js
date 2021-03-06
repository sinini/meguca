/*
Main redis controller module
 */

// Main database controller class
class Yakusoku extends events.EventEmitter {
	constructor(board, ident) {
		super();
		this.id = ++(cache.YAKUMAN);
		this.board = board;

		//Should moderation be allowed on this board?
		this.isContainmentBoard	= config.containment_boards.indexOf(board) > -1;
		this.ident = ident;
		this.subs = [];
	}
	target_key(id) {
		return id === 'live' ? 'board:' + this.board : 'thread:' + id;
	}
	check_throttle(ip, callback) {
		// So we can spam new threads in debug mode
		if (config.DEBUG)
			return callback();
		redis.exists(`ip:${ip}:throttle:thread`, (err, exists) => {
			if (err)
				return callback(err);
			callback(exists && Muggle('Too soon.'));
		});
	}
	get_banner(cb) {
		redis.get('banner:info', cb);
	}
	set_banner(message, cb) {
		redis.set('banner:info', message, err => {
			if (err)
				return cb(err);

			// Dispatch new banner
			const m = redis.multi();
			this._log(m, 0, common.UPDATE_BANNER, [message]);
			m.exec(cb);
		});
	}
}

exports.Yakusoku = Yakusoku;

function postKey(num, op) {
	return `${op == num ? 'thread' : 'post'}:${num}`;
}
