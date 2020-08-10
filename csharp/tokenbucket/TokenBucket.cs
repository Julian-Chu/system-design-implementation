using System;

namespace tokenbucket
{
    public class TokenBucket
    {
        private readonly long _maxBucketSize;
        private readonly long _refillRate;

        private long _currentBucketSize;
        private long _lastRefillTimestamp;

        public TokenBucket(long maxBucketSize, long refillRate)
        {
            _maxBucketSize = maxBucketSize;
            _refillRate = refillRate;
        }

        public bool AllowRequest(int tokens)
        {
            Refill();
            
            if (_currentBucketSize > tokens)
            {
                _currentBucketSize -= tokens;
                return true;
            }

            return false;
        }

        private void Refill()
        {
            var now = DateTime.Now.Date.Ticks;
            var  tokensToAdd = (now - _lastRefillTimestamp) * _refillRate / 1e9;
            _currentBucketSize = Math.Min(_currentBucketSize + (long)tokensToAdd, _maxBucketSize);
            _lastRefillTimestamp = now;
        }
    }
}