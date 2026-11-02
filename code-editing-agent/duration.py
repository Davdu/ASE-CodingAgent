import re

_PATTERN = re.compile(r"(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?")


def to_seconds(text):
    """Convert a duration like "1h30m", "45s" or "2h" to seconds.

    Raises ValueError if the text is empty or not a valid duration.
    """
    if not isinstance(text, str) or not text:
        raise ValueError(f"invalid duration: {text!r}")
    match = _PATTERN.fullmatch(text)
    if match is None:
        raise ValueError(f"invalid duration: {text!r}")
    hours, minutes, seconds = (int(g) if g else 0 for g in match.groups())
    return hours * 3600 + minutes * 60 + seconds
